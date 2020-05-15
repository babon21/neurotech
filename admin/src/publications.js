import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm, NumberInput } from 'react-admin';
import { FileInput, FileField } from 'react-admin';


export const PublicationList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="year" sortable={false}/>
            <TextField source="title" sortable={false}/>
            <FileField source="file.url" title="file.name" />
        </Datagrid>
    </List>
);

const PublicationTitle = ({ record }) => {
    return <span>Publication {record ? `"${record.title}"` : ''}</span>;
};

export const PublicationEdit = props => (
    <Edit title={<PublicationTitle />} {...props}>
        <SimpleForm>
            <TextInput disabled source="id" />
            <NumberInput source="year" />
            <TextInput multiline source="title" />
            <FileInput source="file" label="Файл публикации">
                <FileField source="url" title="name" />
            </FileInput>
        </SimpleForm>
    </Edit>
);

export const PublicationCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <NumberInput source="year" />
            <TextInput multiline source="title" />
            <FileInput source="file" label="Файл публикации">
                <FileField source="url" title="name" />
            </FileInput>
        </SimpleForm>
    </Create>
);
