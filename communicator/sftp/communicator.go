package sftp

import (
	sshComm "github.com/mitchellh/packer/communicator/ssh"
	"github.com/mitchellh/packer/packer"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
)

type comm struct {
	config  *sshComm.Config
	ssh     packer.Communicator
	address string
}

// Creates a new packer.Communicator implementation, using SFTP. This takes
// an already existing TCP connection and SSH configuration.
func New(address string, config *sshComm.Config) (result *comm, err error) {
	sshConn, err := sshComm.New(address, config)
	if err != nil {
		return
	}

	result = &comm{
		config:  config,
		ssh:     sshConn,
		address: address,
	}
	return
}

func (c *comm) Start(cmd *packer.RemoteCmd) (err error) {
	err = c.ssh.Start(cmd)
	return
}

func (c *comm) Upload(path string, input io.Reader, fi *os.FileInfo) error {
	conn, err := c.config.Connection()
	if err != nil {
		return err
	}

	log.Printf("handshaking with SSH")
	sshConn, sshChan, req, err := ssh.NewClientConn(conn, c.address, c.config.SSHConfig)
	if err != nil {
		log.Printf("handshake error: %s", err)
		return err
	}

	sshClient := ssh.NewClient(sshConn, sshChan, req)
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer client.Close()

	f, err := client.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, input); err != nil {
		return err
	}

	return nil
}

func (c *comm) UploadDir(dst string, src string, excl []string) error {
	panic("unimplemented")
}

func (c *comm) Download(path string, output io.Writer) error {
	panic("unimplemented")
}
