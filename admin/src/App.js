import React from 'react';
import { Admin, Resource } from 'react-admin';
import simpleRestProvider from 'ra-data-simple-rest';
import { NewsList, NewsEdit, NewsCreate } from './news';
import { PublicationList, PublicationEdit, PublicationCreate } from './publications';
import { StudentWorkList, StudentWorkEdit, StudentWorkCreate } from './studentWorks';
import { DisciplinesList, DisciplineEdit, DisciplineCreate } from './disciplines';
import authProvider from './authProvider';


const dataProvider = simpleRestProvider('http://localhost:8080');
const App = () => (
  <Admin dataProvider={dataProvider} authProvider={authProvider}>
    <Resource name="news" list={NewsList} edit={NewsEdit} create={NewsCreate} />
    <Resource name="publications" list={PublicationList} edit={PublicationEdit} create={PublicationCreate} />
    <Resource name="student-works" list={StudentWorkList} edit={StudentWorkEdit} create={StudentWorkCreate} />
    <Resource name="disciplines" list={DisciplinesList} edit={DisciplineEdit} create={DisciplineCreate} />
  </Admin>
);
export default App;
