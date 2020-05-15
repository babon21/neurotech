import simpleRestProvider from 'ra-data-simple-rest';


const dataProvider = simpleRestProvider('http://localhost:8080');

const customDataProvider = {
    ...dataProvider,
    update: (resource, params) => {
        if (resource === 'disciplines') {
            return updateDiscipline(resource, params)
        }
        return dataProvider.update(resource, params);

    },

    create: (resource, params) => {
        if (resource === 'disciplines') {
            return createDiscipline(resource, params)
        }
        return dataProvider.create(resource, params);

    },
};

function createDiscipline(resource, params) {
    if (params.data.hasOwnProperty('lections') === false && params.data.hasOwnProperty('books') === false) {
        return dataProvider.create(resource, params);
    }

    var lections = null
    if (typeof params.data.lections !== 'undefined' && params.data.lections != null) {
        for (let i = 0; i < params.data.lections.length; i++) {
            if (params.data.lections[i].hasOwnProperty('rawFile') === true) {
                lections = params.data.lections.splice(i, params.data.lections.length - i)
            }
        }
    }

    var books = null
    if (typeof params.data.books !== 'undefined' && params.data.books != null) {
        for (let i = 0; i < params.data.books.length; i++) {
            if (params.data.books[i].hasOwnProperty('rawFile') === true) {
                books = params.data.books.splice(i, params.data.books.length - i)
            }
        }
    }

    return dataProvider.create(resource, params).then((result) => {
        if (lections == null && books == null) {
            return new Promise((resolve) => resolve(result));
        }

        let formData = new FormData();
        fillFormWithLections(formData, lections)
        fillFormWithBooks(formData, books)

        return fetch('http://localhost:8080/disciplines/' + result.data.id + '/files', {
            method: 'POST',
            body: formData
        }).then(() => result)
    })
}

function updateDiscipline(resource, params) {
    if (params.data.hasOwnProperty('lections') === false && params.data.hasOwnProperty('books') === false) {
        return dataProvider.update(resource, params);
    }

    var lections = null
    if (typeof params.data.lections !== 'undefined' && params.data.lections != null) {
        for (let i = 0; i < params.data.lections.length; i++) {
            if (params.data.lections[i].hasOwnProperty('rawFile') === true) {
                lections = params.data.lections.splice(i, params.data.lections.length - i)
            }
        }
    }

    var books = null
    if (typeof params.data.books !== 'undefined' && params.data.books != null) {
        for (let i = 0; i < params.data.books.length; i++) {
            if (params.data.books[i].hasOwnProperty('rawFile') === true) {
                books = params.data.books.splice(i, params.data.books.length - i)
            }
        }
    }

    return dataProvider.update(resource, params).then((result) => {
        if (lections == null && books == null) {
            return new Promise((resolve) => resolve(result));
        }

        let formData = new FormData();
        fillFormWithLections(formData, lections)
        fillFormWithBooks(formData, books)

        return fetch('http://localhost:8080/disciplines/' + params.data.id + '/files', {
            method: 'POST',
            body: formData
        }).then(() => result)
    })
}

function fillFormWithLections(formData, lections) {
    if (typeof lections !== 'undefined' && lections !== null) {
        for (let i = 0; i < lections.length; i++) {
            formData.append("lections", lections[i].rawFile)
        }
    }
}

function fillFormWithBooks(formData, books) {
    if (typeof books !== 'undefined' && books !== null) {
        for (let i = 0; i < books.length; i++) {
            formData.append("books", books[i].rawFile)
        }
    }
}

export default customDataProvider;
