import React, { Component, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { longpoll } from "../../redux/Longpoll";
import { Row, Col, List, Space, Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import Contacts from '../layout/Contacts';

const { Search } = Input;
const Home = () => {

    return (
        <Row style={{ height: "100vh" }}>
            <Col span={4}>
                <Contacts />
            </Col>
            <Col span={20}>right</Col>
        </Row >
    )
}


export default Home;