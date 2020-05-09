import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm } from 'react-admin';



export const NewsList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="title" sortable={false}/>
            <TextField source="content" sortable={false}/>
        </Datagrid>
    </List>
);

const NewsTitle = ({ record }) => {
    return <span>News {record ? `"${record.title}"` : ''}</span>;
};

export const NewsEdit = props => (
    <Edit title={<NewsTitle />} {...props}>
        <SimpleForm>
            <TextInput disabled source="id" />
            <TextInput source="title" />
            <TextInput multiline source="content" />
        </SimpleForm>
    </Edit>
);

export const NewsCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="title" />
            <TextInput multiline source="content" />
        </SimpleForm>
    </Create>
);
