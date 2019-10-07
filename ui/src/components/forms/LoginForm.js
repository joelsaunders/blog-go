import React from "react";
import {Field, Form, Formik} from "formik";
import {connect} from "react-redux";
import {signIn} from "../../actions";
import theBookOfJoel from "../../apis/theBookOfJoel";

const LoginForm = (props) => {
    const handleSubmit = async (values, actions) => {
        let response;
        try {
            response = await theBookOfJoel.post(
                'api/v1/user/login',
                values
            );
        } catch (e) {
            if (e.response.status === 400) {
                actions.setErrors(e.response.data);
            }
            actions.setSubmitting(false);
            actions.setStatus({ msg: JSON.stringify(e.response.data)});
            return
        }
        actions.setSubmitting(false);
        props.setAuth(values.email, response.data.token);
        props.onDismiss();
        return response
    };

    return <div>
        <Formik
            initialValues={{email: "", password:""}}
            onSubmit={(values, actions) => handleSubmit(values, actions)}
        >
            { props => {
                const { status, isSubmitting} = props;

                return (
                    <Form>
                        <div className="mb-4">
                            <label
                                className="block text-gray-700 text-sm font-bold mb-2"
                                htmlFor="email"
                            >
                                Email
                            </label>
                            <Field
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                type="text"
                                name="email"
                            />
                        </div>
                        <div className="mb-4">
                            <label
                                className="block text-gray-700 text-sm font-bold mb-2"
                                htmlFor="password"
                            >
                                Password
                            </label>
                            <Field
                                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                                type="password"
                                name="password"
                            />
                            {status && status.msg && <div>{status.msg}</div>}
                        </div>
                        <button
                            className="inline-block text-sm px-4 py-2 leading-none border rounded text-teal-500 border-white hover:border-transparent hover:text-white hover:bg-teal-500 mt-4"
                            type="submit"
                            disabled={isSubmitting}
                        >
                            Login
                        </button>
                    </Form>
                );
            }}
        </Formik>
    </div>
};

export default connect(null, {setAuth: signIn})(LoginForm);