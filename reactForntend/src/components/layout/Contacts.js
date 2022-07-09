import React, { Component, useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { longpoll } from "../../redux/Longpoll";
import { Row, Col, List, Space, Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import agent from "../../agent"

const { Search } = Input;
const Contacts = () => {
    const [loading, setLoading] = useState(false)
    const [searchUser, setSearchUser] = useState([])
    const data = useSelector(state => state.msg.contacts)

    const onSearch = async (value) => {
        setLoading(true)
        try {
            const result = await agent.User.getUserInfo(value)
            setSearchUser(result)
        } catch (e) {
            console.log(e)
        } finally {
            setLoading(false)
        }
    }

    return (
        <List
            style={{ height: "100%" }}
            header={<Search placeholder="input search text" onSearch={onSearch} />}
            bordered
            grid={{ column: 1, }}
            dataSource={[...searchUser, ...data]}
            renderItem={(item) => (
                <div style={{ display: "inline-flex" }}>
                    {item}
                </div>

            )}
        />
    )
}


export default Contacts;