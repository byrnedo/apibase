package dockertest
import (
	"os"
	gDoc "github.com/fsouza/go-dockerclient"
)
var (
	dockCli *gDoc.Client
)


func init() {
	var err error
	dockCli, err = gDoc.NewClientFromEnv()
	if err != nil {
		panic("Failed to connect to docker:" + err.Error())
	}
}


func Start(image string, portB map[gDoc.Port][]gDoc.PortBinding) (string, error) {

	if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: image, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {

		return "", err
	}

	con, err := dockCli.CreateContainer(gDoc.CreateContainerOptions{
		Config: &gDoc.Config{
			Labels: map[string]string{
				"ApiBaseTestFlag" : image,
			},
			Image: image,
		},
		HostConfig: &gDoc.HostConfig{
			PortBindings: portB,
		},
	})
	if err != nil {
		return "", err
	}

	if err := dockCli.StartContainer(con.ID, nil); err != nil {
		return "", err
	}
	return con.ID, nil
}

func Running(image string) (string, error) {
	cons, err :=dockCli.ListContainers(gDoc.ListContainersOptions{
		Filters: map[string][]string{
			"label":[]string{"ApiBaseTestFlag=" + image},
		},
	})
	if err != nil {
		return "", err
	}

	if len(cons) == 0 {
		return "", nil
	}
	return cons[0].ID, nil
}

func Stop(image string) (bool, error) {
	var (
		id string
		err error
	)
	if id, err = Running(image); err != nil && len(id) > 0 {
		return false, err
	}
	if err = dockCli.RemoveContainer(gDoc.RemoveContainerOptions{
		Force: true,
		ID: id,
	}); err != nil {
		return false,err
	}
	return true, nil

}
