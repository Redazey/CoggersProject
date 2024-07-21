import React, {useState, useEffect, useCallback} from "react";
import { Link } from 'react-router-dom';
import {RpcError} from "grpc-web";
import {ServerInfo, getServersInfo} from '../components/gRPC/servParser/grpcServParser'
import serversConfig from '../config/servers.json'

interface ServerCardProps {
    server: ServerInfo;
}

const ServerCardLoad: React.FC<ServerCardProps> = ({ server }) => {
    console.log(server)
    return (
        <div className="servers__card">
            <Link to={`/server/${server.Name}`}><h2>{server.Name}</h2></Link>
            <table>
                <tbody>
                    <tr>
                        <td>Онлайн:</td>
                        <td>{server.Online}/{server.MaxOnline}</td>
                    </tr>
                    <tr>
                        <td>Адрес:</td>
                        <td>{server.Adress}</td>
                    </tr>
                    <tr>
                        <td>Версия:</td>
                        <td>{server.Version}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    );
}

const Aside: React.FC = () => {
    const [serversInfo, setServerInfo] = useState<ServerInfo[] | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);

    const fetchServersInfo = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const ServParserResponse = await getServersInfo();
                setServerInfo(ServParserResponse);
            } catch (error) {
                if (error instanceof RpcError) {
                    if (error.code === 5) {
                        setEnd(true);
                    } else {
                        console.error('Error fetching servers info:', error);
                    }
                } else {
                    console.error('Unexpected error:', error);
                }
            } finally {
                setLoading(false);
            }
        }
    }, [loading, end]);

    useEffect(() => {
        setServerInfo(serversConfig.servers);
        fetchServersInfo();

        const interval = setInterval(() => {
            fetchServersInfo();
        }, 5000);

        // Очистка интервала при размонтировании компонента
        return () => {
            clearInterval(interval);
        };
    }, []);

    if (!serversInfo) {
        return <div>{serversInfo}</div>
    }

    return (
        <>
            {serversInfo.map((value) => (
                <ServerCardLoad server={value}/>
            ))}
        </>
    );
};

export default Aside;
