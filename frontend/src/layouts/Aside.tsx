import React, {useState, useEffect, useCallback} from "react";
import { Link } from 'react-router-dom';
import {RpcError} from "grpc-web";
import {ServerInfo, getServersInfo} from '../components/gRPC/servParser/grpcServParser'

interface ServerCardProps {
    server: ServerInfo;
}

const ServerCard: React.FC<ServerCardProps> = ({ server }) => {
    return (
        <div className="servers__card">
            <Link to={`/server/${server.Name}`}><h2>{server.Name}</h2></Link>
            <table>
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
            </table>
        </div>
    );
}

const Aside: React.FC = () => {
    const [serversInfo, setServerInfo] = useState<Map<string, ServerInfo> | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);

    const fetchServersInfo = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const ServParserResponse = await getServersInfo();
                setServerInfo(ServParserResponse);
                console.log("запуск функции")
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
        fetchServersInfo();

        const interval = setInterval(() => {
            fetchServersInfo();
        }, 60000);

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
            <div className="servers">
                {Array.from(serversInfo).map(([key, value]) => (
                    <ServerCard key={key} server={value}/>
                ))}
            </div>
        </>
    );
};

export default Aside;
