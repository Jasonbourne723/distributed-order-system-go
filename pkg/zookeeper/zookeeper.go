package zookeeper

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkServers = []string{"127.0.0.1:2181"} // Zookeeper服务器地址
	lockPath  = "/lock"                    // 锁的路径
)

type ZookeeperLock struct {
	conn *zk.Conn
}

func NewZookeeperLock() *ZookeeperLock {
	conn, _, err := zk.Connect(zkServers, time.Second*35)
	if err != nil {
		fmt.Println(fmt.Errorf("zookeeper connect failed,%w", err))
	}
	return &ZookeeperLock{conn: conn}
}

// Create a sequential ephemeral node in Zookeeper to represent a lock.
func (zl *ZookeeperLock) AcquireLock(name string) (string, error) {
	// Try to create a sequential, ephemeral node under /lock path
	lockNode := lockPath + "/" + name + "-"
	nodeName, err := zl.conn.Create(lockNode, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", fmt.Errorf("unable to create lock node: %v", err)
	}
	return nodeName, nil
}

// Watch and listen for the previous node in the sequence to be deleted.
func (zl *ZookeeperLock) WatchAndWaitForLock(name string) (string, error) {
	lockNode := lockPath + "/" + name + "-"
	nodeName, err := zl.conn.Create(lockNode, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", fmt.Errorf("unable to create lock node: %v", err)
	}
	for {
		children, _, err := zl.conn.Children(lockPath)
		if err != nil {
			return "", fmt.Errorf("unable to list lock children: %v", err)
		}

		// Sort nodes to determine the predecessor
		var previousNode string
		var currentNodeIdx int

		for idx, child := range children {
			if child == nodeName {
				currentNodeIdx = idx
				break
			}
		}

		// Watch for the predecessor node
		if currentNodeIdx > 0 {
			previousNode = children[currentNodeIdx-1]
		} else {
			// If there's no predecessor, it means this node is the first one, so no need to watch
			return nodeName, nil
		}

		// Watch the previous node for changes
		_, _, ch, err := zl.conn.ExistsW(lockPath + "/" + previousNode)
		if err != nil {
			return "", fmt.Errorf("unable to watch previous node: %v", err)
		}
		// Wait for the previous node to be deleted
		<-ch
	}

}

// Release the lock by deleting the lock node
func (zl *ZookeeperLock) ReleaseLock(nodeName string) error {
	err := zl.conn.Delete(nodeName, -1)
	if err != nil {
		return fmt.Errorf("unable to release lock: %v", err)
	}
	return nil
}
