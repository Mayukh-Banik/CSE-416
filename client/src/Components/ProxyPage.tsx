import React, { useState,useEffect } from 'react';
import Sidebar from "./Sidebar";
import useProxyHostsStyles from '../Stylesheets/ProxyPageStyles';
import { Button, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography, Box, TextField } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface ProxyHost {
  name: string;
  location: string;
  logs: string[];
  Statistics: {
    uptime: string;
  };
  bandwidth: string;
  isEnabled: boolean;
  price: string;
  isHost: boolean;
}


const ProxyHosts: React.FC = () => {
  const [data, setData] = useState<ProxyHost[]>([]); // Corrected to ProxyHost array
  const [input, setInput] = useState<string>('');
  const [proxyHosts, setProxyHosts] = useState<ProxyHost[]>([]); // State to store proxy hosts
  const [currentIP, setCurrentIP] = useState<string>('');
  const [connectedProxy, setConnectedProxy] = useState<ProxyHost | null>(null);
  const [proxyHistory, setProxyHistory] = useState<{ name: string; location: string; timestamp: string }[]>([]);
  const [showHistoryOnly, setShowHistoryOnly] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);
  const [newProxy, setNewProxy] = useState<ProxyHost>({
    name: '',
    location: '',
    logs: [],
    Statistics: { uptime: '' },
    bandwidth: '',
    isEnabled: false,
    price: '',
    isHost: false
  });
  const styles = useProxyHostsStyles();
  const theme = useTheme();
  
  const fetchData = async () => {
    try {
      const response = await fetch(`http://localhost:8081/proxy-data/`);
      if (!response.ok) {
        throw new Error('Failed to fetch data');
      }
      const result = await response.json();
      console.log(result)
      setProxyHosts(result ||[]);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };
  const sendData = async () => {
    try {
      const response = await fetch('http://localhost:8081/proxy-data/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newProxy),
      });
  
      if (!response.ok) {
        throw new Error('Failed to add proxy');
      }
  
      const result = await response.json();
      console.log('Proxy added successfully:', result);
    } catch (error) {
      console.error('Error adding proxy:', error);
      alert('There was an error adding the proxy. Please try again.');
    }
  };
  
  useEffect(() => {
    // Fetch public IP using an API like ipify or ipinfo
    const fetchIP = async () => {
      try {
        const response = await fetch('https://api.ipify.org?format=json');
        const data = await response.json();
        setCurrentIP(data.ip); // Update the state with the user's public IP
      } catch (error) {
        console.error("Error fetching IP:", error);
      }
    }; 
    fetchIP()
    fetchData();
  }, []);

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

  const handleAddProxy = async () => {
    if (newProxy.location.trim() === '' || newProxy.price.trim() === '' || newProxy.Statistics.uptime === '' || newProxy.bandwidth.trim() === '') {
      alert('Please fill in all fields.');
      return;
    }
    await sendData();

    setProxyHosts([...proxyHosts, { ...newProxy, logs: [], isEnabled: false }]);

    // Reset new proxy fields
    setNewProxy({
      name:'',
      location: '',
      logs: [],
      Statistics: { uptime: '' },
      bandwidth: '',
      isEnabled: false,
      price: '',
      isHost:false,
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
                      label="Name"
                      variant="outlined"
                      value={newProxy.name}
                      onChange={(e) => setNewProxy({ ...newProxy, name: e.target.value })}
                    />
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
                      value={newProxy.Statistics.uptime}
                      onChange={(e) => setNewProxy({ ...newProxy, Statistics: { ...newProxy.Statistics, uptime: e.target.value } })}
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
                        <TableCell>{host.location}</TableCell>
                        <TableCell>{host.price}</TableCell>
                        <TableCell>{host.Statistics && host.Statistics.uptime ? host.Statistics.uptime : 'N/A'}</TableCell>
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