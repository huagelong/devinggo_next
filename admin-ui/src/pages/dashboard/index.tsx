import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic } from 'antd';
import { UserOutlined, FileTextOutlined, EyeOutlined } from '@ant-design/icons';
import { Line } from '@ant-design/plots';
import request from '../../utils/request';

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<any>({});
  const [chartData, setChartData] = useState<any[]>([]);

  useEffect(() => {
    // 模拟数据接口调用
    // request.get('/dashboard/statistics').then(setStats);
    // request.get('/dashboard/loginChart').then(setChartData);
    
    // 占位模拟数据
    setStats({
      users: 120,
      activeUsers: 15,
      views: 1024,
    });

    setChartData([
      { date: '2026-03-01', logins: 30 },
      { date: '2026-03-02', logins: 45 },
      { date: '2026-03-03', logins: 28 },
      { date: '2026-03-04', logins: 50 },
      { date: '2026-03-05', logins: 65 },
      { date: '2026-03-06', logins: 55 },
      { date: '2026-03-07', logins: 80 },
    ]);
  }, []);

  const config = {
    data: chartData,
    xField: 'date',
    yField: 'logins',
    point: { size: 5, shape: 'diamond' },
    tooltip: { showMarkers: true },
    state: { active: { style: { shadowBlur: 4, stroke: '#000', fill: 'red' } } },
    interactions: [{ type: 'marker-active' }],
  };

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Card>
            <Statistic title="总用户数" value={stats.users} prefix={<UserOutlined />} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="活跃用户" value={stats.activeUsers} prefix={<EyeOutlined />} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="总访问量" value={stats.views} prefix={<FileTextOutlined />} />
          </Card>
        </Col>
      </Row>

      <Card title="近期登录趋势" style={{ marginTop: 24 }}>
        <Line {...config} />
      </Card>
    </div>
  );
};

export default Dashboard;
