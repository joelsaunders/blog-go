import React from 'react';

const teamData = [
    {
        name: "Joel Saunders",
        position: "Back-end Developer",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494528759/12107892_10156064054045526_4654662208648396258_n_nhzpdy.jpg"
    },{
        name: "Joel Saunders",
        position: "Front-end Developer",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494528817/13939456_10157217881185526_4327777202494721174_n_obrdpq.jpg"
    },{
        name: "Joel Saunders",
        position: "Devops Engineer",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494530180/11009107_10155235419370526_7023342185511609798_n_stdgnj.jpg"
    },{
        name: "Joel Saunders",
        position: "Joel Force One Pilot",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494528893/10384751_10155235412420526_8922016018664635108_n_ubjkqk.jpg"
    },{
        name: "Joel Saunders",
        position: "Janitor",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494177669/joelpic1_gobweb.jpg"
    },{
        name: "Joel Saunders",
        position: "Joel Saunders",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494530245/1016877_10154792592115154_4422530397631181048_n_y08ufu.jpg"
    },{
        name: "Joel Saunders",
        position: "Head Chef",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494528843/13510998_1238570646153379_5610182370510387014_n_jsluoq.jpg"
    },{
        name: "Joel Saunders",
        position: "CEO, CTO, CFO, CSV",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494528870/15747390_10155086779899410_1032417764450909787_n_m3yl6b.jpg"
    },{
        name: "Joel Saunders",
        position: "Chief of Meetings",
        image: "https://res.cloudinary.com/dceeo2a79/image/upload/v1494530338/20170511_120316_1024_xwbsaq.jpg"
    },
];


const TeamCard = ({image, name, position}) => {
    return <div className="w-full sm:w-1/2 md:w-1/3 mb-4 sm:mb-8">
        <div className="w-56 mx-auto rounded bg-white shadow-lg pb-2">
            <img src={image}
                 className="rounded-t h-56 w-full object-cover"
            />
            <div className="m-2">
                <p className="text-center text-base text-gray-900 font-bold">
                    {position}
                </p>
                <p className="text-center text-sm text-gray-600">
                    {name}
                </p>
            </div>
        </div>
    </div>

};

const TeamTable = () => {
    return <div className="w-full my-6 flex flex-col sm:flex-row justify-center flex-wrap">
        {teamData.map(({image, position, name}) => {
            return <TeamCard key={image} image={image} position={position} name={name} />
        })}
    </div>;
};

export default TeamTable;
