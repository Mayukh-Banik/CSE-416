// TransactionPage.tsx
import React, {useState} from "react";
import { Box, Container, CssBaseline, TextField } from "@mui/material";
import useTransactionStyles from "../Stylesheets/TransactionPageStyles";
import TransactionTable from "./TransactionTable";
import Sidebar from "./Sidebar";
import Header from "./Header"; 

const TransactionPage: React.FC = () => {
    const classes = useTransactionStyles();
    const [search, setSearch] = useState("");
    // const [filter, setFilter] = useState("all")
    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearch(event.target.value); //filler for now
    };

    return (
        <Box sx={{ display: "flex" }}>
        <CssBaseline />
        <Sidebar /> 
        <Container component="main" sx={{ flexGrow: 1, p: 3 }}>
            {/* <Header />  */}
            <TextField
                label="Search Transactions"
                variant="outlined"
                onChange={handleSearch}
                value={search}
                sx={{ marginBottom: 2 }}
            />
            {/* add searching transactions */}
            <TransactionTable /> 
        </Container>
        </Box>
    );        
};

export default TransactionPage;
