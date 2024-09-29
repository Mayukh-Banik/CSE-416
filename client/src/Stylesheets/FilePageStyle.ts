import { makeStyles } from "@mui/styles";
import { Theme } from "@mui/material/styles"; // Import Theme

const FilePageStyles = makeStyles((theme: Theme) => ({
    fileContent:{
        display: 'flex',
        flexDirection: 'column',
        marginTop: theme.spacing(2), // Adding some space from the header

    }
}));

export default FilePageStyles;
