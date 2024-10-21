import React from "react";
import { useLocation, Link } from "react-router-dom";
import Sidebar from "./Sidebar";
import {
    Box,
    Container,
    Typography,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { useTheme } from '@mui/material/styles';

interface Account {
    id: number;
    name: string;
}

interface File {
    id: number;
    name: string;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const SearchPage: React.FC = () => {
    const theme = useTheme();
    const location = useLocation();
    
    const queryParams = new URLSearchParams(location.search);
    const searchQuery = queryParams.get("q"); 

    const accountsData: Account[] = [
        { id: 1, name: "john_doe" },
        { id: 2, name: "jane_smith" },
        { id: 3, name: "bob_jones" },
    ];

    const filesData: File[] = [
        { id: 1, name: "file1.txt" },
        { id: 2, name: "file2.txt" },
    ];

    const filteredAccounts = accountsData.filter(account =>
        account.name.toLowerCase().includes(searchQuery?.toLowerCase() || '')
    );

    const filteredFiles = filesData.filter(file =>
        file.name.toLowerCase().includes(searchQuery?.toLowerCase() || '')
    );

    return (
        <Box
            sx={{
                padding: 2,
                marginTop: '70px',
                marginLeft: `${drawerWidth}px`, 
                transition: 'margin-left 0.3s ease', 
                [theme.breakpoints.down('sm')]: {
                    marginLeft: `${collapsedDrawerWidth}px`, 
                },
            }}
        >
            <Sidebar />
            <Container>
                <Typography variant="h4">Search Results</Typography>
                {searchQuery ? (
                    <Typography variant="body1">Results for "{searchQuery}"</Typography>
                ) : (
                    <Typography variant="body1">Please enter a search term.</Typography>
                )}

                {/* Accounts Table */}
                <Accordion>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography variant="h6">Accounts</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <TableContainer component={Paper}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Account ID</TableCell>
                                        <TableCell>Account Name</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {filteredAccounts.length > 0 ? (
                                        filteredAccounts.map((account) => (
                                            <TableRow key={account.id}>
                                                <TableCell>{account.id}</TableCell>
                                                <TableCell>
                                                    <Link 
                                                        to={`/account/${account.name}`} 
                                                        style={{ textDecoration: 'none', color: 'blue' }} 
                                                    >
                                                        {account.name}
                                                    </Link>
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    ) : (
                                        <TableRow>
                                            <TableCell colSpan={2} align="center">
                                                No accounts found.
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </AccordionDetails>
                </Accordion>

                <Accordion>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                        <Typography variant="h6">Files</Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <TableContainer component={Paper}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>File ID</TableCell>
                                        <TableCell>File Name</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {filteredFiles.length > 0 ? (
                                        filteredFiles.map((file) => (
                                            <TableRow key={file.id}>
                                                <TableCell>{file.id}</TableCell>
                                                <TableCell>{file.name}</TableCell>
                                            </TableRow>
                                        ))
                                    ) : (
                                        <TableRow>
                                            <TableCell colSpan={2} align="center">
                                                No files found.
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </AccordionDetails>
                </Accordion>
            </Container>
        </Box>
    );
}

export default SearchPage;
