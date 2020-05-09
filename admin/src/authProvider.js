import qs from 'qs';

export default {
    // called when the user attempts to log in
    login: ({ username, password }) => {
        const credentials = qs.stringify({ 'username': username, 'password': password });

        const request = new Request('http://localhost:8080/auth/login', {
            method: 'POST',
            body: credentials,
            headers: new Headers({ 'Content-Type': 'application/x-www-form-urlencoded' }),
        });

        return fetch(request)
            .then(response => {
                console.log(response.status)
                if (response.status !== 200) {
                    console.log("status != 200")
                    throw new Error(response.statusText);
                }

                console.log("response status = 200")

                let type = response.headers.get('Content-Type')
                if (type !== 'text/html; charset=UTF-8') {
                    console.log(type)
                    throw new Error("Username or password is incorrect");
                }

                localStorage.setItem('isAuth', 'some_str');
                return response;
            }, () => {
                throw new Error("Server is not available");
            });
    },
    // called when the user clicks on the logout button
    logout: () => {
        localStorage.removeItem('isAuth');
        return Promise.resolve();
    },
    // called when the API returns an error
    checkError: ({ status }) => {
        if (status === 401 || status === 403) {
            localStorage.removeItem('isAuth');
            return Promise.reject();
        }
        return Promise.resolve();
    },
    // called when the user navigates to a new location, to check for authentication
    checkAuth: () => {
        return localStorage.getItem('isAuth')
            ? Promise.resolve()
            : Promise.reject();
    },
    // called when the user navigates to a new location, to check for permissions / roles
    getPermissions: () => Promise.resolve(),
};
