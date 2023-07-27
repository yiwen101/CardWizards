// src/App.js

import React from 'react';
import { Admin, Resource} from 'react-admin';
import dataProvider from './dataProvider';
import { ServiceList, ServiceCreate, ServiceEdit } from './components/ServiceList';
import { ApiList, ApiEdit } from './components/APIList';
import { RouteList, RouteEdit, RouteCreate } from './components/RouteList';
import DnsIcon from '@mui/icons-material/Dns';
import ApiIcon from '@mui/icons-material/Api';
import RouteIcon from '@mui/icons-material/Route';



const App = ({addr}) => (
  <Admin dataProvider={dataProvider(addr)}>
    <Resource name="service" list={ServiceList} edit = {ServiceEdit} create = {ServiceCreate} icon = {DnsIcon}/>
    <Resource name="api" list={ApiList} edit={ApiEdit} hasCreate={false} icon={ApiIcon}/>
    <Resource name="route" list={RouteList} edit={RouteEdit} create = {RouteCreate}  icon={RouteIcon}/>
  </Admin>
);

export default App;
