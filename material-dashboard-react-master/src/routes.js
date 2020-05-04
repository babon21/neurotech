// @material-ui/icons
import LibraryBooks from "@material-ui/icons/LibraryBooks";
import MenuBookIcon from '@material-ui/icons/MenuBook';
import GroupIcon from '@material-ui/icons/Group';
// core components/views for Admin layout
import NewsPage from "views/News.js";
import StudentPage from "views/Students.js";
import StudyMaterialsPage from "views/StudyMaterials.js";
import PublicationPage from "views/Publications.js";


const dashboardRoutes = [
    {
        path: "/news",
        name: "Новости",
        icon: GroupIcon,
        component: NewsPage,
        layout: "/admin"
    },
    {
        path: "/publications",
        name: "Публикации",
        icon: LibraryBooks,
        component: PublicationPage,
        layout: "/admin"
    },
    {
        path: "/study-materials",
        name: "Учебные материалы",
        icon: MenuBookIcon,
        component: StudyMaterialsPage,
        layout: "/admin"
    },
    {
        path: "/students",
        name: "Студенты",
        icon: GroupIcon,
        component: StudentPage,
        layout: "/admin"
    }
];

export default dashboardRoutes;