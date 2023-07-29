import React from 'react';
import {Datagrid, List, TextField, TextInput,Edit, SimpleForm, EditButton, Toolbar, useRedirect, SaveButton, DeleteButton, Create} from 'react-admin';
import Tooltip from '@material-ui/core/Tooltip';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';



const ServiceFilter = [
  <TextInput source="serviceName" label="Service Name" alwaysOn/>,
  <TextInput source="methodName" label="Method Name" alwaysOn/>,
];
export const RouteList = () => (
  <List   filters = {ServiceFilter} >
        <Datagrid >
            <TextField source="serviceName" />
            <TextField source="methodName" />
            <TextField source="httpMethod" />
            <TextField source="url" />
            <EditButton basePath='/route' />
            <DeleteButton basePath='/route' />
        </Datagrid>
      </List>
);

export const RouteEdit = (props) => (
    <Edit {...props}>
        <SimpleForm toolbar={<RouteEditToolbar />} >
            <TextInput source="id" disabled/>
            <TextInput source="serviceName" disabled />
            <TextInput source="methodName" disabled/>
            <TextInput source="httpMethod" />
            <TextInput source="url" />
        </SimpleForm>
    </Edit>
);

export const RouteCreate = (props) => (
  <Create {...props}>
    <SimpleForm toolbar = {<RoutePostToolbar/>}>
    <TextInput source="serviceName"  />
    <TextInput source="methodName"  />
    <TextInput source="httpMethod" />
    <TextInput source="url" />
    </SimpleForm>
  </Create>
);

const RouteEditToolbar = props => {
  const redirect = useRedirect();

  const handleReturn = () => {
    redirect('/route');
  };

  return (
    <Toolbar {...props} >
       <SaveButton />

      {/* Custom Return button */}
      <Tooltip title="Return to the Route List" aria-label="return">
        <div onClick={handleReturn} style={{ cursor: 'pointer' }}>
          <ArrowBackIcon />
        </div>
      </Tooltip>
      <DeleteButton basePath='/route' />
    </Toolbar>
  );
};

const RoutePostToolbar = props => {
  const redirect = useRedirect();

  const handleReturn = () => {
    redirect('/route');
  };

  return (
    <Toolbar {...props}>
      <SaveButton />
      <Tooltip title="Return to the Service List" aria-label="return">
        <div onClick={handleReturn} style={{ cursor: 'pointer' }}>
          <ArrowBackIcon />
        </div>
      </Tooltip>
    </Toolbar>
  );
}






