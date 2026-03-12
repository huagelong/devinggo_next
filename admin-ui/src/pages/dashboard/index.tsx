import React, { useState, useEffect } from 'react';
import { Row, Col, Card } from 'antd';
import { ProCard, StatisticCard } from '@ant-design/pro-components';
import { Line, Column } from '@ant-design/plots';

const { Statistic } = StatisticCard;

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<any>({});
  const [chartData, setChartData] = useState<any[]>([]);

  useEffect(() => {
    // 占位模拟数据
    setStats({
      users: 1256,
      activeUsers: 843,
      views: 45024,
      revenue: 124500,
    });

    setChartData([
      { date: '2026-03-01', logins: 130 },
      { date: '2026-03-02', logins: 145 },
      { date: '2026-03-03', logins: 128 },
      { date: '2026-03-04', logins: 150 },
      { date: '2026-03-05', logins: 165 },
      { date: '2026-03-06', logins: 155 },
      { date: '2026-03-07', logins: 180 },
    ]);
  }, []);

  const config = {
    data: chartData,
    xField: 'date',
    yField: 'logins',
    smooth: true,
    area: {
      style: {
        fill: 'l(270) 0:#ffffff 0.5:#7ec2f3 1:#1890ff',
      },
    },
  };

  return (
    <div className="space-y-6">
      <ProCard ghost gutter={[16, 16]}>
        <ProCard colSpan={6} layout="center" bordered>
          <StatisticCard
            statistic={{
              title: '总用户数',
              value: stats.users,
              description: <Statistic title="较上周" value="8.04%" trend="up" />,
            }}
          />
        </ProCard>
        <ProCard colSpan={6} layout="center" bordered>
          <StatisticCard
            statistic={{
              title: '活跃用户',
              value: stats.activeUsers,
              description: <Statistic title="较上周" value="3.14%" trend="down" />,
            }}
          />
        </ProCard>
        <ProCard colSpan={6} layout="center" bordered>
          <StatisticCard
            statistic={{
              title: '总访问量',
              value: stats.views,
              description: <Statistic title="较昨日" value="12.4%" trend="up" />,
            }}
          />
        </ProCard>
        <ProCard colSpan={6} layout="center" bordered>
          <StatisticCard
            statistic={{
              title: '总营收',
              value: stats.revenue,
              prefix: '¥',
              description: <Statistic title="较上月" value="5.8%" trend="up" />,
            }}
          />
        </ProCard>
      </ProCard>

      <Row gutter={16}>
        <Col span={16}>
          <ProCard title="近期访问趋势" bordered headerBordered>
            <div className="h-[350px]">
              <Line {...config} />
            </div>
          </ProCard>
        </Col>
        <Col span={8}>
          <ProCard title="快捷操作" bordered headerBordered>
            <div className="grid grid-cols-3 gap-4 mb-4">
              {['系统设置', '权限管理', '用户管理', '角色管理', '部门分配', '操作日志'].map(item => (
                <div key={item} className="flex flex-col items-center justify-center p-4 bg-gray-50 rounded-lg cursor-pointer hover:bg-blue-50 hover:text-blue-600 transition-colors">
                  <span className="text-sm">{item}</span>
                </div>
              ))}
            </div>
          </ProCard>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
