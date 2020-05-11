import simpleRestProvider from 'ra-data-simple-rest';


const dataProvider = simpleRestProvider('http://localhost:8080');

const customDataProvider = {
    ...dataProvider,
    update: (resource, params) => {
        if (resource !== 'disciplines') {
            return dataProvider.update(resource, params);
        }
    
        var files = null
        for (let i = 0; i < params.data.files.length; i++) {
            if (params.data.files[i].hasOwnProperty('rawFile') === true) {
                files = params.data.files.splice(i, params.data.files.length - i)
            }
        }

        return dataProvider.update(resource, params).then((result) => {
            if (files == null) {
                return new Promise((resolve) => resolve(result));
            }

            const formData = new FormData();
            for (let i = 0; i < files.length; i++) {
                formData.append("multiplefiles", files[i].rawFile)
            }

            return fetch('http://localhost:8080/disciplines/' + params.data.id + '/files', {
                method: 'POST',
                body: formData
            }).then(() => result)
        })
    },

    create: (resource, params) => {
        if (resource !== 'disciplines') {
            return dataProvider.create(resource, params);
        }       

        if (params.data.hasOwnProperty('files') === false) {
            return dataProvider.create(resource, params);
        }

        var files = null
        for (let i = 0; i < params.data.files.length; i++) {
            if (params.data.files[i].hasOwnProperty('rawFile') === true) {
                files = params.data.files.splice(i, params.data.files.length - i)
            }
        }

        return dataProvider.create(resource, params).then((result) => {
            if (files == null) {
                return new Promise((resolve) => resolve(result));
            }

            // var files = params.data.files
            const formData = new FormData();
            for (let i = 0; i < files.length; i++) {
                formData.append("multiplefiles", files[i].rawFile)
            }

            return fetch('http://localhost:8080/disciplines/' + result.data.id + '/files', {
                method: 'POST',
                body: formData
            }).then(() => result)
        })
    },
};

export default customDataProvider;
