import React, { Component } from 'react';
import { login } from '../../../redux/Auth'
import { connect } from 'react-redux'
import { Button, Space, Form, Input, Row, Col } from 'antd';
import { LockOutlined, UserOutlined } from '@ant-design/icons';

const mapStateToProps = (dispatch) => {
    return {
        dispatch,
    }
}

class Login extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            loading: false
        };
        this.login = this.login.bind(this)
    }



    async login(v) {
        this.setState({ loading: true })
        await this.props.dispatch(login({ username: v.username, password: v.password }))
        this.setState({ loading: false })
    }


    render() {
        return (
            <Row justify="center" align="middle" style={{ height: "100vh" }}>
                <Form
                    style={{ maxWidth: "300px" }}
                    initialValues={{ remember: true }}
                    onFinish={this.login}
                    // onFinishFailed={onFinishFailed}
                    autoComplete="off"
                >
                    <Form.Item name="username" rules={[
                        {
                            required: true,
                            message: 'Please input your Username!',
                        },
                    ]}>
                        <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" />
                    </Form.Item>
                    <Form.Item name="password" rules={[
                        {
                            required: true,
                            message: 'Please input your Password!',
                        },
                    ]}>
                        <Input.Password prefix={<LockOutlined className="site-form-item-icon" />} placeholder="Password" />
                    </Form.Item>
                    <Button type="primary" htmlType="submit" loading={this.state.loading}>
                        Submit
                    </Button>
                </Form>
            </Row >
        );
    }
}

export default connect(mapStateToProps /** no second argument */)(Login);