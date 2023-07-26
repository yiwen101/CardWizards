import { BooleanField, Datagrid, List, TextField, TextInput,BooleanInput,Edit, SimpleForm, EditButton} from 'react-admin';


const ServiceFilter = [
    <TextInput source="serviceName" label="Service Name" alwaysOn/>,
  ];


export const ApiList = () => (
    <List >
        <Datagrid >
            <TextField source="id"/>
            <BooleanField source="isSleeping" />
            <BooleanField source="validationOn" />
            <EditButton basePath='/api' />
        </Datagrid>
    </List>
);

export const ApiEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <TextInput source="id" disabled/>
            <BooleanInput source="isSleeping"/>
            <BooleanInput source="validationOn"/>
        </SimpleForm>
    </Edit>
);