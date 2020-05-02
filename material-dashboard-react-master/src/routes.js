// @material-ui/icons
import Dashboard from "@material-ui/icons/Dashboard";
import Person from "@material-ui/icons/Person";
import LibraryBooks from "@material-ui/icons/LibraryBooks";
import BubbleChart from "@material-ui/icons/BubbleChart";
import LocationOn from "@material-ui/icons/LocationOn";
import Notifications from "@material-ui/icons/Notifications";
import Unarchive from "@material-ui/icons/Unarchive";
import Language from "@material-ui/icons/Language";
import MenuBookIcon from '@material-ui/icons/MenuBook';
import GroupIcon from '@material-ui/icons/Group';
// core components/views for Admin layout
import DashboardPage from "views/Dashboard/Dashboard.js";
import UserProfile from "views/UserProfile/UserProfile.js";
import TableList from "views/TableList/TableList.js";
import Typography from "views/Typography/Typography.js";
import NotificationsPage from "views/Notifications/Notifications.js";


const dashboardRoutes = [{
        path: "/dashboard",
        name: "Dashboard",
        rtlName: "لوحة القيادة",
        icon: Dashboard,
        component: DashboardPage,
        layout: "/admin"
    },
    {
        path: "/user",
        name: "User Profile",
        rtlName: "ملف تعريفي للمستخدم",
        icon: Person,
        component: UserProfile,
        layout: "/admin"
    },
    {
        path: "/table",
        name: "Table List",
        rtlName: "قائمة الجدول",
        icon: "content_paste",
        component: TableList,
        layout: "/admin"
    },
    {
        path: "/typography",
        name: "Typography",
        rtlName: "طباعة",
        icon: LibraryBooks,
        component: Typography,
        layout: "/admin"
    },
    {
        path: "/notifications",
        name: "Notifications",
        rtlName: "إخطارات",
        icon: Notifications,
        component: NotificationsPage,
        layout: "/admin"
    },
    {
        path: "/educational-materials",
        name: "Educational materials",
        rtlName: "إخطارات",
        icon: MenuBookIcon,
        component: NotificationsPage,
        layout: "/admin"
    },
    {
        path: "/students",
        name: "Students",
        rtlName: "إخطارات",
        icon: GroupIcon,
        component: NotificationsPage,
        layout: "/admin"
    }
];

export default dashboardRoutes;