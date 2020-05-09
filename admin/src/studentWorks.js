import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm, NumberInput } from 'react-admin';


export const StudentWorkList = props => (
    <List {...props}>
        <Datagrid rowClick="edit">
            <TextField source="student" sortable={false} />
            <TextField source="year" sortable={false} />
            <TextField source="title" sortable={false} />
        </Datagrid>
    </List>
);

const StudentWorkTitle = ({ record }) => {
    return <span>Student Work {record ? `"${record.title}"` : ''}</span>;
};

export const StudentWorkEdit = props => (
    <Edit title={<StudentWorkTitle />} {...props}>
        <SimpleForm>
            <TextInput disabled source="id" />
            <TextInput source="student" />
            <NumberInput source="year" />
            <TextInput source="title" />
        </SimpleForm>
    </Edit>
);

export const StudentWorkCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="student" />
            <NumberInput source="year" />
            <TextInput source="title" />
        </SimpleForm>
    </Create>
);
