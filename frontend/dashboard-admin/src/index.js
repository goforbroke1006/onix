import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import reportWebVitals from './reportWebVitals';

import {BrowserRouter, Link, Route, Routes} from "react-router-dom";
import MainPage from './pages/MainPage';
import SourceListPage from "./pages/SourceListPage";
import CriteriaListPage from './pages/CriteriaListPage';
import ApiDataProvider from "./external/api-data-provider";

const apiDataProvider = new ApiDataProvider();

ReactDOM.render(
    <BrowserRouter>
        <div>

            <ul>
                <li><Link to={"/"}>Main page</Link></li>
                <li><Link to={"/source/list"}>Source list</Link></li>
                <li><Link to={"/criteria/list"}>Criteria list</Link></li>
            </ul>

            <Routes>
                <Route path={"/"} element={<MainPage provider={apiDataProvider}/>} exact/>
                <Route path={"/source/list"} element={<SourceListPage provider={apiDataProvider}/>}/>
                <Route path={"/criteria/list"} element={<CriteriaListPage provider={apiDataProvider}/>}/>
            </Routes>

        </div>
    </BrowserRouter>



    /*<React.StrictMode>
      <App />
    </React.StrictMode>*/,
    document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
