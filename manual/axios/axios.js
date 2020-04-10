const getBtn = document.getElementById('get-btn');
const postBtn = document.getElementById('post-btn');
const putBtn = document.getElementById('put-btn');
const deleteBtn = document.getElementById('delete-btn');

const getData = () => {
    axios.get('http://127.0.0.1:8080/axios').then(response => {
        console.log(response);
    });
};

const sendData = () => {
    axios
        .post(
            'http://127.0.0.1:8080/axios', {
                email: 'eve.holt@reqres.in'
                    // password: 'pistol'
            }, {
                // headers: {
                //   'Content-Type': 'application/json'
                // }
            }
        )
        .then(response => {
            console.log(response);
        })
        .catch(err => {
            console.log(err, err.response);
        });
};

const deleteData = () => {
    axios
        .delete("http://127.0.0.1:8080/axios")
        .then(response => {
            console.log(response);
        })
        .catch(err => {
            console.log(err, err.response);
        });
}

const putData = () => {
    axios
        .put("http://127.0.0.1:8080/axios")
        .then(response => {
            console.log(response);
        })
        .catch(err => {
            console.log(err, err.response);
        });
}

getBtn.addEventListener('click', getData);
postBtn.addEventListener('click', sendData);
putBtn.addEventListener('click', putData);
deleteBtn.addEventListener('click', deleteData);