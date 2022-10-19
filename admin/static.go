package admin

import "net/http"

func (admin *adminPanel) handleIndex(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title></title>
</head>
<body>
Admin Panel
</body>
</html>
		`))
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}
