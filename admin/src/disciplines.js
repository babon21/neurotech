import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm } from 'react-admin';
import { FileInput, FileField } from 'react-admin';
import { BooleanField } from 'react-admin';
import { BooleanInput } from 'react-admin';

const DisciplinesPanel = ({ id, record, resource }) => {
    console.log(record.lections)
    return (
        record.lections.map((file) => {
            return <p><a href={file.url}>{file.name}</a></p>
        })
    );
}

export const DisciplinesList = props => (
    <List {...props}>
        <Datagrid rowClick="edit" expand={<DisciplinesPanel />}>
            <TextField source="name" sortable={false} />
            <BooleanField source="is_current_semester" />
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
            <BooleanInput label="Текущий семестр?" source="is_current_semester" />
            <FileInput source="lections" label="Лекции" multiple={true}>
                <FileField source="url" title="name" />
            </FileInput>

            <FileInput source="books" label="Учебные и методические пособия" multiple={true}>
                <FileField source="url" title="name" />
            </FileInput>
        </SimpleForm>
    </Edit>
);

export const DisciplineCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <BooleanInput label="Текущий семестр?" source="is_current_semester  " />
            <FileInput source="lections" label="Лекции" multiple={true}>
                <FileField source="src" title="name" />
            </FileInput>

            <FileInput source="books" label="Учебные и методические пособия" multiple={true}>
                <FileField source="books" title="name" />
            </FileInput>
        </SimpleForm>
    </Create>
);
