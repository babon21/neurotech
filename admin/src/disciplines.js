import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm } from 'react-admin';


export const DisciplinesList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="name" sortable={false}/>
        </Datagrid>
    </List>
);

const DisciplinesTitle = ({ record }) => {
    return <span>Discipline {record ? `"${record.title}"` : ''}</span>;
};

export const DisciplineEdit = props => (
    <Edit title={<DisciplinesTitle />} {...props}>
        <SimpleForm>
            <TextInput disabled source="id" />
            <TextInput source="name" />
        </SimpleForm>
    </Edit>
);

export const DisciplineCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="title" />
            <TextInput source="name" />
        </SimpleForm>
    </Create>
);
