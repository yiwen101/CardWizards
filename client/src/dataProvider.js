// to do, allow modification
const apiUrl = 'http://127.0.0.1:8080/admin'; // Replace with your actual backend API URL

const buildID = (resource) => (item) => {
    return resource =="route"
        ? ({...item, id:item.httpMethod+"/"+item.url})
        : resource =="api"
        ? ({...item, id:item.serviceName+"/"+item.methodName})
        : {...item, id:item.serviceName};
}

const dataProvider = {
    getList: async (resource, params) => {
        const url = `${apiUrl}/${resource}`;
        const response = await fetch(url);
        const json = await response.json();
        // why need bracker outside{}
        let j = json.map(buildID(resource));
        return {
            data: j,
            total: j.length,
        };
    },

    getOne: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const response = await fetch(url);
        const json = await response.json();
        return {
            data: buildID(resource)(json),
        };
    },

    create: async (resource, params) => {
        const url = `${apiUrl}/${resource}`;
        const options = {
            method: 'POST',
            body: JSON.stringify(params.data),
        };
        const response = await fetch(url, options);
        const json = await response.json();
        return {
            data: { ...params.data, id: json.id },
        };
    },
    update: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const options = {
            method: 'PUT',
            body: JSON.stringify(params.data),
        };
        const response = await fetch(url, options);
        const json = await response.json();
        return {
            data: json,
        };
    },

    delete: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const options = {
            method: 'DELETE',
        };
        await fetch(url, options);
        return {
            data: params.previousData,
        };
    }
};

export default dataProvider;
