import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm } from 'react-admin';
import { DateInput, DateField } from 'react-admin';


export const NewsList = props => (
    <List {...props}>
        <Datagrid rowClick="show">
            <TextField source="title" sortable={false}/>
            <DateField source="publication_date" sortable={false}/>
            <TextField source="content" sortable={false}/>
        </Datagrid>
    </List>
);

const NewsTitle = ({ record }) => {
    return <span>News {record ? `"${record.title}"` : ''}</span>;
};

export const NewsEdit = props => { return(
    <Edit title={<NewsTitle />} {...props}>
        <SimpleForm>
            <TextInput disabled source="id" />
            <TextInput source="title" />
            <DateInput source="publication_date" defaultValue={publicationDate}/>
            <TextInput multiline source="content" />
        </SimpleForm>
    </Edit>
);};

const publicationDate = new Date();
export const NewsCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="title" />
            <DateInput source="publication_date" defaultValue={publicationDate}/>
            <TextInput multiline source="content" />
        </SimpleForm>
    </Create>
);
