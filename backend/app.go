package backend

import "fmt"

func Run() error {
	root := findRootDir()
	if root == "" {
		return fmt.Errorf("couldnt find root dir")
	}

    err := ensureChildDirsExist(root); 
	if err != nil {
        return err
    }

    err = ensureEssentialFilesExist(root); 
	if err != nil {
        return err
    }

    configureServer()

    return runServer(":8080")
}
