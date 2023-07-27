import React, {useState} from 'react'
import{List, Datagrid, TextField, EditButton, DeleteButton, BooleanField,  Create, SimpleForm, TextInput, BooleanInput, TopToolbar, FilterButton, CreateButton, Edit, useRedirect, Toolbar, SaveButton,  } from 'react-admin'
import HelpIcon from '@material-ui/icons/Help';
import Tooltip from '@material-ui/core/Tooltip';
import { makeStyles } from '@material-ui/core/styles';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
//import ServiceFilter from './utils';

const ListActions = () => {
  return (
  <TopToolbar>
    <FilterButton label="Filter" alwaysOn />
    <CreateButton basePath='/service' alwaysOn />
  </TopToolbar>
);}

const ServiceFilter = [
  <TextInput source="serviceName" label="Service Name"  alwaysOn/>,
];

export function ServiceList(props) {
  return (
    <List  exporter={false} hasCreate={true} filters = {ServiceFilter} {...props}>
      <Datagrid>
        <TextField source='serviceName' label="Service Name" />
        <TextField source='clusterName' />
        <TextField source='loadBalanceOption' />
        <BooleanField source="isSleeping" label="Sleeping" />
        <EditButton basePath='/service' />
        <DeleteButton basePath='/service' />
      </Datagrid>
    </List>
  )
}

export const ServiceCreate = (props) => (
  <Create {...props}>
    <SimpleForm toolbar = {<ServicePostToolbar/>}>
      <TextInput source="idlFileName" label="IDL FileName" />
      <TextInput source="clusterName" />
      <FieldWithTooltip label="" tooltip="'weighted random' or 'weighted round robin'">
        <TextInput source="loadBalanceOption" label="load balance option" />
      </FieldWithTooltip>
      <BooleanInput source="isSleeping" label="Is Service Sleeping" />
    </SimpleForm>
  </Create>
);

const useStyles = makeStyles(theme => ({
  tooltip: {
    fontSize: theme.typography.pxToRem(14),
  },
}));

const FieldWithTooltip = ({ label, tooltip, children }) => {
  const classes = useStyles();

  return (
    <div style={{ display: 'flex', alignItems: 'center', flexDirection:'row-reverse' }}>
      <div style={{ marginRight: 5 }}>
        <Tooltip title={tooltip} aria-label="tooltip" classes={{ tooltip: classes.tooltip }}>
          <HelpIcon fontSize="small" />
        </Tooltip>
      </div>
      <div style={{flex:1}}>
        {children}
        <div style={{ marginLeft: 22, color: '#999' }}>{label}</div>
      </div>
    </div>
  );
};

const ServiceEditToolbar = props => {
  const redirect = useRedirect();

  const handleReturn = () => {
    redirect('/service');
  };

  return (
    <Toolbar {...props}>
      <SaveButton />
      <Tooltip title="Return to the Service List" aria-label="return">
        <div onClick={handleReturn} style={{ cursor: 'pointer' }}>
          <ArrowBackIcon />
        </div>
      </Tooltip>
      <DeleteButton />
    </Toolbar>
  );
};

const ServicePostToolbar = props => {
  const redirect = useRedirect();

  const handleReturn = () => {
    redirect('/service');
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

export const ServiceEdit = props => {
  const [reloadIdl, setReloadIdl] = useState(false);

  const handleReloadIdlChange = (event) => {
    setReloadIdl(event.target.checked);
  };

  return (<Edit {...props}>
    <SimpleForm toolbar={<ServiceEditToolbar />}>
    <TextInput source="id" disabled/>
  
      <FieldWithTooltip label="" tooltip="'weighted random' or 'weighted round robin'">
        <TextInput source="loadBalanceOption" label="load balance option" />
      </FieldWithTooltip>
        <TextInput source="clusterName" label="cluster name" />
      <FieldWithTooltip label="" tooltip="Leave this field blank unless there is change to the idl file, in which case need to re-upload by inputing the full file name with extensions">
        <TextInput source="idlFileName" disabled = {!reloadIdl} />
      </FieldWithTooltip>
      <BooleanInput source="isSleeping" label="Is Sleeping" />
      <BooleanInput source="reloadIdl" label="Reload IDL" onChange={handleReloadIdlChange} />
    </SimpleForm>
  </Edit>
);}


