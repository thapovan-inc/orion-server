package authprovider

func GetNameSpaceFromAuthToken(authToken string) (string, error) {
	var namespace = "default" //TODO: Check if namespace is stored in the config and then assign the namespace
	if authToken != "" {
		// TODO: Extract namespace from authToken
	}
	return namespace, nil
}
