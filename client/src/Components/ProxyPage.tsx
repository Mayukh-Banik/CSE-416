import React, { useState } from 'react';
import Sidebar from "./Sidebar";
import useProxyHostsStyles from '../Stylesheets/ProxyPageStyles';
import { Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Box, TextField, Paper } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface ProxyHost {
  access: string;
  location: string;
  logs: string[];
  statistics: {
    uptime: string; // E.g., "99.9%"
  };
  bandwidth: string; // E.g., "100 Mbps"
  isEnabled: boolean;
  price: string;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100; 

const ProxyHosts: React.FC = () => {
  const styles = useProxyHostsStyles();
  const theme = useTheme();
  const [proxyHosts, setProxyHosts] = useState<ProxyHost[]>([
    {
      access: 'Public',
      location: 'New York, USA',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: {
        uptime: '99.9%',
      },
      bandwidth: '100 Mbps',
      isEnabled: false,
      price: 'Free',
    },
    {
      access: 'Public',
      location: 'London, UK',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: {
        uptime: '98.7%',
      },
      bandwidth: '200 Mbps', 
      isEnabled: false,
      price: 'Free',
    },
    {
      access: 'Private',
      location: 'Berlin, Germany',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: {
        uptime: '95.5%',
      },
      bandwidth: '150 Mbps', 
      isEnabled: false,
      price: '$20',
    },
  ]);

  const [currentIP, setCurrentIP] = useState<string>('192.168.0.1');
  const [connectedProxy, setConnectedProxy] = useState<ProxyHost | null>(null);

  const [newProxy, setNewProxy] = useState<ProxyHost>({
    access: 'Private',
    location: '',
    logs: [],
    statistics: {
      uptime: '',
    },
    bandwidth: '', 
    isEnabled: false,
    price: '',
  });

  const [showForm, setShowForm] = useState<boolean>(false); // State to control the visibility of the form

  const handleConnect = (host: ProxyHost) => {
    const updatedHosts = proxyHosts.map(h => ({
      ...h,
      isEnabled: h.location === host.location ? true : h.isEnabled,
    }));

    setProxyHosts(updatedHosts);
    setConnectedProxy(host);
    alert(`Connected to ${host.location}`);
  };

  const handleAddProxy = () => {
    if (newProxy.location.trim() === '' || newProxy.price.trim() === '' || newProxy.statistics.uptime === '' || newProxy.bandwidth.trim() === '') {
      alert('Please fill in all fields.');
      return;
    }

    setProxyHosts([...proxyHosts, { ...newProxy, logs: [], isEnabled: false }]);
    
    // Reset new proxy fields
    setNewProxy({
      access: 'Private',
      location: '',
      logs: [],
      statistics: {
        uptime: '',
      },
      bandwidth: '', // Reset bandwidth
      isEnabled: false,
      price: '',
    });

    // Hide the form after adding
    setShowForm(false);
  };

  const handleSortByLocation = () => {
    const sortedHosts = [...proxyHosts].sort((a, b) => a.location.localeCompare(b.location));
    setProxyHosts(sortedHosts);
  };

  const handleSortByPrice = () => {
    const sortedHosts = [...proxyHosts].sort((a, b) => {
      const priceA = a.price === 'Free' ? 0 : parseFloat(a.price.replace('$', ''));
      const priceB = b.price === 'Free' ? 0 : parseFloat(b.price.replace('$', ''));
      return priceA - priceB;
    });
    setProxyHosts(sortedHosts);
  };

  // Function to handle viewing statistics
  const handleViewStatistics = (stats: { uptime: string }) => {
    alert(`Uptime: ${stats.uptime}`);
  };

  return (
    <div className={styles.container}>
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
      <Box sx={{ marginTop: 2 }}>
        <Typography variant="h4">Proxy Hosts</Typography>
        <Typography variant="h6">Your Current IP: {currentIP}</Typography>  
        <br />
        
        {/* Add Yourself as Proxy Button and Sort Buttons on the same row */}
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 1 }}>
          <Button variant="contained" onClick={() => setShowForm(prev => !prev)}>
            {showForm ? 'Hide Form' : 'Add Yourself as Proxy'}
          </Button>
          {/* Sort Buttons */}
          <Box sx={{ display: 'flex', gap: '10px' }}>
            <Button variant="outlined" onClick={handleSortByLocation}>
              Sort by Location
            </Button>
            <Button variant="outlined" onClick={handleSortByPrice}>
              Sort by Price
            </Button>
          </Box>
        </Box>

        {/* Expandable Form Section */}
        {showForm && (
          <Box sx={{ marginTop: 1 }}>
            <Typography variant="h6">Fill in your proxy details</Typography>
            <Box sx={{ display: 'flex', gap: 2 }}>
              <TextField
                label="Location"
                variant="outlined"
                value={newProxy.location}
                onChange={(e) => setNewProxy({ ...newProxy, location: e.target.value })}
              />
              <TextField
                label="Price ($)"
                variant="outlined"
                type="number"
                value={newProxy.price}
                onChange={(e) => setNewProxy({ ...newProxy, price: e.target.value })}
                InputProps={{ inputProps: { min: 0 } }}
              />
              <TextField
                label="Uptime (%)"
                variant="outlined"
                value={newProxy.statistics.uptime}
                onChange={(e) => setNewProxy({ ...newProxy, statistics: { ...newProxy.statistics, uptime: e.target.value } })}
              />
              <TextField
                label="Bandwidth"
                variant="outlined"
                value={newProxy.bandwidth}
                onChange={(e) => setNewProxy({ ...newProxy, bandwidth: e.target.value })}
              />
              <Button variant="contained" onClick={handleAddProxy}>
                Add Proxy
              </Button>
            </Box>
          </Box>
        )}

        {connectedProxy && (
          <Typography variant="h6" sx={{ marginTop: 2 }}>
            Connected Proxy: {connectedProxy.location}
          </Typography>
        )}

        <TableContainer>
          <Table className={styles.table} component={Paper} sx={{marginTop:2}}>
            <TableHead>
              <TableRow>
                <TableCell>Access</TableCell>
                <TableCell>Location</TableCell>
                <TableCell>Price</TableCell>
                <TableCell>Uptime</TableCell>
                <TableCell>Bandwidth</TableCell> 
                <TableCell>Logs</TableCell>
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {proxyHosts.map((host, index) => (
                <TableRow key={index}>
                  <TableCell>{host.access}</TableCell>
                  <TableCell>{host.location}</TableCell>
                  <TableCell>{host.price}</TableCell>
                  <TableCell>{host.statistics.uptime}</TableCell>
                  <TableCell>{host.bandwidth}</TableCell> {/* Displaying bandwidth */}
                  <TableCell>
                    <Button variant="text" onClick={() => alert(host.logs.join('\n'))}>
                      View Logs
                    </Button>
                  </TableCell>
                  <TableCell>
                    <Button variant="contained" onClick={() => handleConnect(host)}>
                      Connect
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
      </Box>
    </div>
  );
};

export default ProxyHosts;
