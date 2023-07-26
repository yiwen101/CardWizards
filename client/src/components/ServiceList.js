import React from 'react'
import{InfiniteList, Datagrid, TextField, EditButton, DeleteButton, BooleanField, ShowButton, Create, SimpleForm, TextInput, BooleanInput, TopToolbar, FilterButton, CreateButton } from 'react-admin'

const ListActions = () => {
  <TopToolbar>
    <FilterButton label="Filter" alwaysOn />
    <CreateButton basePath='/service' alwaysOn />
  </TopToolbar>
}

const ServiceFilter = [
  <TextInput source="serviceName" label="Service Name" alwaysOn/>,
];

export function ServiceList(props) {
  return (
    <InfiniteList  exporter={false} hasCreate={true} filters = {ServiceFilter} {...props}>
      <Datagrid>
        <TextField source='serviceName' label="Service Name" />
        <TextField source='clusterName' />
        <TextField source='loadBalanceOption' />
        <BooleanField source="isSleeping" label="Sleeping" />
        <ShowButton label = "APIs" />
        <EditButton basePath='/service' />
        <DeleteButton basePath='/service' />
      </Datagrid>
    </InfiniteList>
  )
}

export const ServiceCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="serviceName" label="Service Name" />
      <TextInput source="clusterName" />
      <TextInput source="loadBalanceOption" />
      <BooleanInput source="isSleeping" label="Is Service Sleeping" />
    </SimpleForm>
  </Create>
);

