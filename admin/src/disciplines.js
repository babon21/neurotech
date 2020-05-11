import React from 'react';
import { List, Datagrid, TextField, Create, Edit, TextInput, SimpleForm } from 'react-admin';
import { FileInput, FileField } from 'react-admin';
import { SelectInput } from 'react-admin';
import { ChipField } from 'react-admin';


const DisciplinesPanel = ({ id, record, resource }) => {
    console.log(record.files)
    return (
        record.files.map((file) => {
            return <p><a href={file.url}>{file.name}</a></p>
        })
    );
}

export const DisciplinesList = props => (
    <List {...props}>
        <Datagrid rowClick="edit" expand={<DisciplinesPanel />}>
            <TextField source="name" sortable={false} />
            <ChipField source="tag" sortable={false} />
            {/* <FileField source="files" src="url" title="name" /> */}
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
            <FileInput source="files" label="Related files" multiple={true}>
                <FileField source="url" title="name" />
            </FileInput>
        </SimpleForm>
    </Edit>
);

export const DisciplineCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <SelectInput source="tag" choices={[
                { id: 'none', name: 'None' },
                { id: 'current_semester', name: 'Текущий семестр' },
            ]} />
            <FileInput source="files" label="Related files" multiple={true}>
                <FileField source="src" title="title" />
            </FileInput>
        </SimpleForm>
    </Create>
);
