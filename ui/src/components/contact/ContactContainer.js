import React from 'react';


const LogoBox = ({link, children}) => {
    return <div className="w-1/2 flex justify-center">
        <a href={link}
           className="my-10 sm:my-auto"
        >
            <svg viewBox="0 0 24 24" width="40" height="40"
                 stroke="currentColor" strokeWidth="2" fill="none"
                 strokeLinecap="round" strokeLinejoin="round"
                 className="hover:text-teal-500 my-10 sm:my-auto"
            >
                {children}
            </svg>
        </a>
    </div>
};


const ContactContainer = () => {
    return <div className="w-full my-6 rounded overflow-hidden shadow-lg bg-white flex flex-col sm:flex-row">
        <div className="w-full sm:w-1/2 px-6 py-4">
            <h3 className="text-gray-900 font-bold text-xl mb-2">Get in touch</h3>
            <p className="text-gray-700 text-base">
                Please feel free to contact me with any questions or queries you might have about topics in this blog, or anything else
            </p>
            <p className="text-gray-700 text-base mt-2">
                I'll try and respond as fast as I can.
            </p>
            <p className="text-gray-700 text-base mt-2">
                Reach me at <a className="hlink" href="mailto:joel.st.saunders@gmail.com">joel.st.saunders@gmail.com</a>
            </p>
        </div>
        <div className="flex flex-row flex flex-wrap w-full sm:w-1/2">
            <LogoBox link="https://www.linkedin.com/in/joel-saunders-266476102">
                    <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/>
                    <rect x="2" y="9" width="4" height="12"/>
                    <circle cx="4" cy="4" r="2"/>
            </LogoBox>
            <LogoBox link="https://www.github.com/joelsaunders">
                    <path
                        d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/>
            </LogoBox>
            <LogoBox link="mailto:joel.st.saunders@gmail.com">
                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                    <polyline points="22,6 12,13 2,6"/>
            </LogoBox>
            <LogoBox link="https://www.instagram.com/hashtagjoel">
                    <rect x="2" y="2" width="20" height="20" rx="5" ry="5"/>
                    <path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/>
                    <line x1="17.5" y1="6.5" x2="17.5" y2="6.5"/>
            </LogoBox>
        </div>
    </div>;
};

export default ContactContainer;
