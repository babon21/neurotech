import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm, RichTextField } from 'react-admin';


export const NewsList = props => (
    <List {...props}>
        <Datagrid rowClick="show">
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
