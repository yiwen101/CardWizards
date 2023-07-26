// src/App.js

import React from 'react';
import { Admin, Resource, ListGuesser, EditGuesser} from 'react-admin';
import dataProvider from './dataProvider';
import { ServiceList, ServiceCreate } from './components/ServiceList';
import { ApiList, ApiEdit } from './components/APIList';
import { RouteList } from './components/RouteList';


const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource name="service" list={ServiceList} create = {ServiceCreate} />
    <Resource name="api" list={ApiList} edit={ApiEdit} />
    <Resource name="route" list={RouteList} />

  </Admin>
);

export default App;
