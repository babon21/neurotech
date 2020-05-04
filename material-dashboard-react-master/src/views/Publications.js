/*eslint-disable*/
import React from "react";
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import PublicationTable from "components/PublicationTable.js";
import ActionTable from "components/TableWithActions.js";

var bugs = [
  'Sign contract for "What are conference organizers afraid of?"',
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  "Create 4 Invisible User Experiences you Never Knew About"
];

export default function Students() {

  return (
    <GridContainer>
      <GridItem xs={12} sm={12} md={12}>
        <Button color="primary">Добавить</Button>
      </GridItem>
      <GridItem xs={12} sm={12} md={12}>
        <PublicationTable
          headerColor="primary"
          tables={[
            {
              tabContent: (
                <ActionTable
                  tableHeaderColor="primary"
                  tableHead={["Год", "Наименование"]}
                  // checkedIndexes={[0, 3]}
                  // tasksIndexes={[0, 1, 2, 3]}
                  tableData={[
                    ["2020", "Практические аспекты применения методов, алгоритмов и средств индуктивного анализа данных в приоритетных отраслях : отчет о НИР / Новосиб. гос. техн. ун-т ; исполн.: А. А. Якименко, О. К. Альсова, А. В. Гаврилов, С. П. Ильиных, О. В. Казанская, А. А. Малявко, В. К. Мищенко, П. В. Мищенко, В. Г. Токарев, Г. В. Трошина ; рук. Е. В. Рабинович. - Новосибирск, 2020. - 216 с. - №АААА-Б20-220012890096-3."],
                    ["2018", "Curaçao"],
                    ["1998", "Netherlands"],
                    ["1998", "Korea, South"],
                    ["2015", "Malawi"],
                    ["2005", "Gloucester"]
                  ]}
                />
              )
            }
          ]}
        />
      </GridItem>
    </GridContainer>
  );
}
