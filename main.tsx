
import * as React from "react";
import * as ReactDOM from "react-dom";

import {Names, Sprites} from './gen'

interface Result {
    Name: string;
    Cost: number;
    Hearts: number;
    Duration: number;
    EffectString: string;
    StrengthString: string;
    Ingredients: string;
}

interface Test {
    Expected: Result;
    Actual: Result;
}

interface TestState {
    Tests: Test[];
    Shown?: Test;
}

// 'HelloProps' describes the shape of props.
// State is never set so we use the 'undefined' type.
 class TestSuite extends React.Component<undefined, TestState> {
    constructor(){
        super();
        this.state = {Tests: []};
    }
    public async componentDidMount(){
        var resp = await fetch("https://raw.githubusercontent.com/captncraig/wildchef/master/results.json");
        var dat = (await resp.json()) as Result[];
        var tsts: Test[] = [];
        for(var r of dat){
            tsts.push({Expected: r, Actual: predict(r)})
        }
        this.setState({Tests: tsts});
    }
    onEnter = (tst: Test) =>{
        this.setState({Shown: tst});
    }
    onLeave = () =>{
        this.setState({Shown: null});
    }
    render() {
        return <div>
            <Detail Test={this.state.Shown} />
            <div style={{marginTop: "200px"}}>
                <h3>{this.state.Tests.length ? <span>{this.state.Tests.length} Test Cases</span>: <span>Loading</span>}</h3>
                {this.state.Tests.map((tst,i) => <SingleTest key={i} Test={tst} OnEnter={this.onEnter.bind(this,tst)} OnLeave={this.onLeave}/>)}
            </div>
        </div>;
    }
}
interface SingleTestProps {
    Test: Test;
    OnEnter(): void;
    OnLeave(): void;
}

class SingleTest extends React.Component<SingleTestProps, undefined>{
    render(){
        var ex = this.props.Test.Expected;
        var style = {
            background: spriteStyle(ex.Name),
        }
        return <div className='test' style={style} onMouseEnter={this.props.OnEnter} onMouseLeave={this.props.OnLeave}>{ex.Name}</div>
    }
}

interface DetailProps {
    Test?: Test;
}


class Detail extends React.Component<DetailProps, undefined>{
    render(){
         var style = {
            position: "fixed" as "fixed",
            top: "0",
            width: '100%',
            height: '200px',
            overflow: "hidden" as "hidden",
            backgroundColor: 'white',
        }
        
        var content: JSX.Element = null
        if (this.props.Test){
            var ex = this.props.Test.Expected;
            var ac = this.props.Test.Actual;
            content = 
<div>
    {ex.Name}
    {ex.Ingredients.split(",").map(ing => {
        var style = {
            background: spriteStyle(ing),
        }
        return <div className='test' style={style}/>
    })}
</div>
        }
        return <div style={style}>{content}</div>
    }
}

ReactDOM.render(
    <TestSuite />,
    document.getElementById("example")
);

function predict(r: Result): Result {
    return {
        Name: Names.DubiousFood,
        Hearts: 4,
        Cost: 2,
    } as Result
}

function spriteStyle(name: string):string{
    var s = Sprites[name] || {X: 7, Y: 31};
    var xoff = s.X * 96;
    var yoff = s.Y * 96;
    return `url(sprites.png) -${xoff}px -${yoff}px, darkgrey`
}

// 0-6. 
function matchPct(a: Result, b: Result): number{
    var matches = 0
    if (a.Name == b.Name){
        matches++;
    }
    if (a.Cost == b.Cost){
        matches++;
    }
    if (a.Duration == b.Duration){
        matches++;
    }
    if (a.EffectString == b.EffectString){
        matches++
    }
    if (a.StrengthString == b.StrengthString){
        matches++
    }
    if (a.Hearts == b.Hearts){
        matches++
    }
    return matches;
}

