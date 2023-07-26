import { Datagrid, List, TextField, UrlField } from 'react-admin';

export const RouteList = () => (
    <List>
        <Datagrid >
            <TextField source="serviceName" />
            <TextField source="methodName" />
            <TextField source="httpMethod" />
            <UrlField source="url" />
        </Datagrid>
    </List>
);