import socket
import threading
import logging
import select
import argparse

# Configure logging
logging.basicConfig(level=logging.INFO, 
                    format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class ProxyForwarder:
    def __init__(self, local_host='0.0.0.0', local_port=9999, 
                 remote_host='192.168.1.100', remote_port=8888):
        """
        Initialize the proxy forwarder.
        
        :param local_host: Local host to bind (0.0.0.0 for network-wide access)
        :param local_port: Local port to listen on
        :param remote_host: Remote proxy server host
        :param remote_port: Remote proxy server port
        """
        self.local_host = local_host
        self.local_port = local_port
        self.remote_host = remote_host
        self.remote_port = remote_port
        self.server_socket = None

    def start(self):
        """
        Start the proxy forwarder and listen for incoming connections.
        """
        try:
            # Create local server socket
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            
            # Bind to local address
            self.server_socket.bind((self.local_host, self.local_port))
            self.server_socket.listen(100)
            
            logger.info(f"Proxy forwarder started on {self.local_host}:{self.local_port}")
            logger.info(f"Forwarding to {self.remote_host}:{self.remote_port}")
            
            while True:
                # Accept incoming client connection
                client_socket, client_address = self.server_socket.accept()
                logger.info(f"Received connection from {client_address}")
                
                # Create a new thread to handle the client
                client_thread = threading.Thread(
                    target=self.handle_client, 
                    args=(client_socket,)
                )
                client_thread.daemon = True
                client_thread.start()
        
        except Exception as e:
            logger.error(f"Error starting proxy forwarder: {e}")
        finally:
            if self.server_socket:
                self.server_socket.close()

    def handle_client(self, client_socket):
        """
        Handle client connection and forward to remote proxy.
        
        :param client_socket: Socket connection from local client
        """
        try:
            # Connect to remote proxy
            remote_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            remote_socket.connect((self.remote_host, self.remote_port))
            
            logger.info(f"Connected to remote proxy {self.remote_host}:{self.remote_port}")
            
            # Tunnel bidirectional traffic
            self.tunnel(client_socket, remote_socket)
        
        except Exception as e:
            logger.error(f"Error handling client: {e}")
        finally:
            client_socket.close()

    def tunnel(self, client_socket, remote_socket):
        """
        Tunnel bidirectional traffic between client and remote proxy.
        
        :param client_socket: Local client-side socket
        :param remote_socket: Remote proxy server socket
        """
        def forward(source, destination):
            try:
                while True:
                    readable, _, _ = select.select([source], [], [], 5)
                    if not readable:
                        break
                    
                    data = source.recv(4096)
                    if not data:
                        break
                    
                    destination.sendall(data)
            except Exception as e:
                logger.error(f"Tunneling error: {e}")
        
        # Create threads for bidirectional forwarding
        client_to_remote = threading.Thread(
            target=forward, 
            args=(client_socket, remote_socket)
        )
        remote_to_client = threading.Thread(
            target=forward, 
            args=(remote_socket, client_socket)
        )
        
        client_to_remote.daemon = True
        remote_to_client.daemon = True
        
        client_to_remote.start()
        remote_to_client.start()
        
        # Wait for threads to complete
        client_to_remote.join()
        remote_to_client.join()
        
        # Close sockets
        client_socket.close()
        remote_socket.close()

def main():
    """
    Create and start the proxy forwarder with command-line arguments.
    """
    # Set up argument parsing
    parser = argparse.ArgumentParser(description='Proxy Forwarder')
    parser.add_argument('-lh', '--local-host', 
                        default='0.0.0.0', 
                        help='Local host to bind (default: 0.0.0.0)')
    parser.add_argument('-lp', '--local-port', 
                        type=int, 
                        default=19484, 
                        help='Local port to listen on (default: 19484)')
    parser.add_argument('-rh', '--remote-host', 
                        required=True,
                        help='Remote proxy server host IP')
    parser.add_argument('-rp', '--remote-port', 
                        type=int, 
                        default=19483, 
                        help='Remote proxy server port (default: 19483)')
    
    # Parse arguments
    args = parser.parse_args()
    
    # Create and start forwarder
    forwarder = ProxyForwarder(
        local_host=args.local_host,
        local_port=args.local_port,
        remote_host=args.remote_host,
        remote_port=args.remote_port
    )
    forwarder.start()

if __name__ == '__main__':
    main()
