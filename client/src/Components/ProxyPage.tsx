import React, { useState,useEffect } from 'react';
import Sidebar from "./Sidebar";
import useProxyHostsStyles from '../Stylesheets/ProxyPageStyles';
import { Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Box, TextField } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface ProxyHost {
  name: string;
  access: string;
  location: string;
  logs: string[];
  statistics: {
    uptime: string;
  };
  bandwidth: string;
  isEnabled: boolean;
  price: string;
}


const ProxyHosts: React.FC = () => {
  const [data, setData] = useState<string>('');
  const [input, setInput] = useState<string>('');
  const styles = useProxyHostsStyles();
  const theme = useTheme();
  const fetchData = async () => {
    try {
      const response = await fetch('http://localhost:8081/proxy-data');
      if (!response.ok) {
        throw new Error('Failed to fetch data');
      }
      const result = await response.json();
      setData(result.message);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };
  const sendData = async () => {
    try {
      const response = await fetch('http://localhost:8081/proxy-data', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: input }),
      });
      if (!response.ok) {
        throw new Error('Failed to send data');
      }
      alert('Data sent successfully!');
    } catch (error) {
      console.error('Error sending data:', error);
    }
  };
  useEffect(() => {
    fetchData();
  }, []);
  
  const [proxyHosts, setProxyHosts] = useState<ProxyHost[]>([
    {
      name:'p1',
      access: 'Public',
      location: 'New York, USA',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: { uptime: '99.9%' },
      bandwidth: '100 Mbps',
      isEnabled: false,
      price: 'Free',
    },
    {
      name:'p2',
      access: 'Public',
      location: 'London, UK',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: { uptime: '98.7%' },
      bandwidth: '200 Mbps',
      isEnabled: false,
      price: 'Free',
    },
    {
      name:'private1',
      access: 'Private',
      location: 'Berlin, Germany',
      logs: ['Log entry 1', 'Log entry 2'],
      statistics: { uptime: '95.5%' },
      bandwidth: '150 Mbps',
      isEnabled: false,
      price: '20 SQD',
    },
  ]);

  //dummy ip address, tentative to change
  const [currentIP, setCurrentIP] = useState<string>('192.168.0.1');
  const [connectedProxy, setConnectedProxy] = useState<ProxyHost | null>(null);

  // Dummy history
  const [proxyHistory, setProxyHistory] = useState<{ name: string, location: string; timestamp: string }[]>([
    { name:'p1',location: 'New York, USA', timestamp: '2024-10-15 10:00:00' },
    { name:'priavte1',location: 'London, UK', timestamp: '2024-10-14 14:30:00' },
  ]);

  const [showHistoryOnly, setShowHistoryOnly] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);

  const [newProxy, setNewProxy] = useState<ProxyHost>({
    name:'',
    access:'',
    location: '',
    logs: [],
    statistics: { uptime: '' },
    bandwidth: '',
    isEnabled: false,
    price: '',
  });

  const handleConnect = (host: ProxyHost) => {
    const updatedHosts = proxyHosts.map(h => ({
      ...h,
      isEnabled: h.location === host.location ? true : h.isEnabled,
    }));

    setProxyHosts(updatedHosts);
    setConnectedProxy(host);

    // Update proxy history
    const newHistoryEntry = { name: host.name,location: host.location, timestamp: new Date().toLocaleString() };
    setProxyHistory([...proxyHistory, newHistoryEntry]);

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
      name:'',
      access: '',
      location: '',
      logs: [],
      statistics: { uptime: '' },
      bandwidth: '',
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

  const handleClearAndShowHistory = () => {
    setShowHistoryOnly(true); // Show history
  };

  const handleReturn = () => {
    setShowHistoryOnly(false); // Return
  };

  return (
    <div className={styles.container}>
      <Box className={styles.boxContainer}>
        <Sidebar />
        <Box sx={{ marginTop: 2 }}>
          <Typography variant="h4">Proxy</Typography>
          <Box className={styles.header}>
            <Typography variant="h6">Your Current IP: {currentIP}</Typography>
            {/* Show Proxy History Button */}
            {!showHistoryOnly && (
              <Button variant="outlined" onClick={handleClearAndShowHistory}>
                Show Proxy History
              </Button>
            )}
          </Box>
          <br />
          {/* Show the history*/}
          {showHistoryOnly ? (
            <>
              <Box className={styles.historyContainer}>
                <Typography variant="h5">Proxy Connection History</Typography>
                <TableContainer className={styles.historyTable}>
                  <Table>
                    <TableHead>
                      <TableRow>
                      <TableCell>Name</TableCell>
                        <TableCell>Location</TableCell>
                        <TableCell>Connected At</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {proxyHistory.map((entry, index) => (
                        <TableRow key={index}>
                          <TableCell>{entry.name}</TableCell>
                          <TableCell>{entry.location}</TableCell>
                          <TableCell>{entry.timestamp}</TableCell>
                          <TableCell>
                            <Button
                              variant="contained"
                              onClick={() => {
                                const host = proxyHosts.find(h => h.location === entry.location);
                                if (host) handleConnect(host);
                              }}
                             >
                              Connect
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
              <Box sx={{ marginTop: 2 }}>
                <Button variant="contained" onClick={handleReturn}>Return to Proxy Hosts</Button>
              </Box>
            </>
          ) : (
            <>
              <Box className={styles.proxyButton}>
                {/* Add Yourself as Proxy  */}
                <Button variant="contained" onClick={() => setShowForm(prev => !prev)}>
                  {showForm ? 'Hide Form' : 'Add Yourself as Proxy'}
                </Button>

                {/* Sort Buttons */}
                <Box sx={{ display: 'flex', gap: '10px' }}>
                  <Button variant="outlined" onClick={handleSortByLocation}>Sort by Location</Button>
                  <Button variant="outlined" onClick={handleSortByPrice}>Sort by Price</Button>
                </Box>
              </Box>

              {/* Expandable Form Section */}
              {showForm && (
                <Box sx={{ marginTop: 1 }} className={styles.form}>
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
                    <Button variant="contained" className={styles.submitButton} onClick={handleAddProxy}>Add Proxy</Button>
                  </Box>
                </Box>
              )}

              <TableContainer className={styles.proxyTable}>
                <Table className={styles.table}>
                  <TableHead>
                    <TableRow>
                    <TableCell>Name</TableCell>
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
                        <TableCell>{host.name}</TableCell>
                        <TableCell>{host.access}</TableCell>
                        <TableCell>{host.location}</TableCell>
                        <TableCell>{host.price}</TableCell>
                        <TableCell>{host.statistics.uptime}</TableCell>
                        <TableCell>{host.bandwidth}</TableCell>
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
            </>
          )}
        </Box>
      </Box>
    </div>
  );
};

export default ProxyHosts;