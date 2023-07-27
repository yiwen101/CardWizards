import {fetchUtils} from "react-admin";
export const dataProvider = (addr) => { 
    const apiUrl = addr + "/admin"
    return {    
    getList: async (resource, params) => {
        const { filter } = params;
        const url = `${apiUrl}/${resource}`;
        const response = await fetchUtils.fetchJson(url);
        const json = await response.json;

        // Check if both "serviceName" and "methodName" filters are present
        if (filter && filter.serviceName && filter.methodName) {
            const { serviceName, methodName } = filter;
            // Filter the data based on both "serviceName" and "methodName"
            const filteredData = json.filter(
                item =>
                    item.serviceName.toLowerCase().includes(serviceName.toLowerCase()) &&
                    item.methodName.toLowerCase().includes(methodName.toLowerCase())
            );
            return {
                data: filteredData,
                total: filteredData.length,
            };
        }

        // Check if there's a filter on the "serviceName" field
        if (filter && filter.serviceName) {
            const { serviceName } = filter;
            // Filter the data based on the "serviceName" field
            const filteredData = json.filter(item =>
                item.serviceName.toLowerCase().includes(serviceName.toLowerCase())
            );
            return {
                data: filteredData,
                total: filteredData.length,
            };
        }

        // Check if there's a filter on the "methodName" field
        if (filter && filter.methodName) {
            const { methodName } = filter;
            // Filter the data based on the "methodName" field
            const filteredData = json.filter(item =>
                item.methodName.toLowerCase().includes(methodName.toLowerCase())
            );
            return {
                data: filteredData,
                total: filteredData.length,
            };
        }

        return {
            data: json,
            total: json.length,
        };
    },


    getOne: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const response = await fetchUtils.fetchJson(url);
        const json = await response.json;
        return {
            data:json,
        };
    },

    create: async (resource, params) => {
        const url = `${apiUrl}/${resource}`;
        const options = {
            method: 'POST',
            body: JSON.stringify(params.data),
        };
        const response = await fetchUtils.fetchJson(url, options);
        const json = await response.json;
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
        const response = await fetchUtils.fetchJson(url, options);
        const json = await response.json;
        return {
            data: json,
        };
    },

    delete: async (resource, params) => {
        const url = `${apiUrl}/${resource}/${params.id}`;
        const options = {
            method: 'DELETE',
        };
        await fetchUtils.fetchJson(url, options);
        return {
            data: params.previousData,
        };
    },
    deleteMany: async (resource, params) => {
        const { ids } = params;
        const deletePromises = ids.map((id) => {
          const url = `${apiUrl}/${resource}/${id}`;
          const options = {
            method: 'DELETE',
          };
          return fetchUtils.fetchJson(url, options);
        });

        await Promise.all(deletePromises);
    
        return {
          data: ids.map((id) => ({ id })),
        };
      },
};
}

export default dataProvider;
