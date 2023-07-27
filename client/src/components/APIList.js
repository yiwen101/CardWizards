import React from 'react';
import { BooleanField, Datagrid, List, TextField, TextInput,BooleanInput,Edit, SimpleForm, EditButton, Toolbar, useRedirect, SaveButton} from 'react-admin';
import Tooltip from '@material-ui/core/Tooltip';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';


const ServiceFilter = [
    <TextInput source="serviceName" label="Service Name" alwaysOn/>,
    <TextInput source="methodName" label="Method Name" alwaysOn/>,
  ];


export const ApiList = () => (
    <List  exporter={false} hasCreate={false} filters = {ServiceFilter} >
        <Datagrid bulkActionButtons={false}>
            <TextField source="serviceName" />
            <TextField source="methodName" />
            <BooleanField source="isSleeping" />
            <BooleanField source="validationOn" />
            <EditButton basePath='/api' />
        </Datagrid>
    </List>
);
  
  export const ApiEdit = props => (
    <Edit {...props}>
      <SimpleForm toolbar={<ApiEditToolbar />}>
        <TextInput source="id" label="ID" disabled />
        <TextInput source="serviceName" label="Service Name" disabled />
        <TextInput source="methodName" label="Method Name" disabled />
        <BooleanInput source="isSleeping" label="Is Sleeping" />
        <BooleanInput source="validationOn" label="Validation On" />
      </SimpleForm>
    </Edit>
  );

const ApiEditToolbar = props => {
    const redirect = useRedirect();
  
    const handleReturn = () => {
      redirect('/api');
    };
  
    return (
      <Toolbar {...props}>
        <SaveButton />
        <Tooltip title="Return to the API List" aria-label="return">
          <div onClick={handleReturn} style={{ cursor: 'pointer' }}>
            <ArrowBackIcon />
          </div>
        </Tooltip>
      </Toolbar>
    );
  };
  
  


