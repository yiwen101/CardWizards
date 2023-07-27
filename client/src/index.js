import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const apiGatewayEndpoint =  process.env.Gateway_Address || 'http://127.0.0.1:8080';

document.title = 'API Gateway Console'; // Replace 'Custom Page Title' with your desired page title

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App addr={apiGatewayEndpoint}/>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
