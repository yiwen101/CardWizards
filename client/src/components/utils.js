// src/components/ReturnButton.js

import React from 'react';
import { useRedirect, TextInput } from 'react-admin';

export const ReturnButton = ({ source }) => {
  const redirect = useRedirect();

  const handleClick = () => {
    redirect(source);
  };

  return <button onClick={handleClick}>Return</button>;
};

export const ServiceFilter = [
  <TextInput source="serviceName" label="Service Name"  alwaysOn/>,
];

export const ServiceAndMethodFilter = [
  <TextInput source="serviceName" label="Service Name"  alwaysOn/>,
  <TextInput source="methodName" label="Method Name"  alwaysOn/>,
]
