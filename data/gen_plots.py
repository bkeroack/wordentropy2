#!/usr/bin/env python

import os
import csv

DATA_PATH = "./stats"

TYPENAME_MAP = {
	"particle": "Plural Article",
	"sarticle": "Singular Article",
	"pnoun": "Plural Noun",
	"snoun": "Singular Noun"
}

import plotly.plotly as py
from plotly.graph_objs import *

with open("plotly_creds.txt", "r") as f:
	creds = [l for l in f]
assert len(creds) == 2

py.sign_in(str(creds[0]).strip(), str(creds[1]).strip())

import numpy as np

type_map = {}
plot_url_map = {}

for filename in os.listdir(DATA_PATH):
	if filename[-4:] == ".csv":
		basename = filename[:-4]
		proper_name = TYPENAME_MAP[basename] if basename in TYPENAME_MAP else str(basename).capitalize()
		with open(os.path.join(DATA_PATH, filename), 'r') as f:
			csvr = csv.reader(f)
			rows = [r for r in csvr]
		bar = Bar(
		        x=[r[0] for r in rows],
		        y=[r[1] for r in rows],
		        name="{} length (characters)".format(proper_name)
		    )
		idata = Data([
		    bar
		])
		layout = Layout(
			title="{} Length Distribution".format(proper_name if basename != "ALL" else "Overall"),
			xaxis=XAxis(title="Word length (characters)"),
			yaxis=YAxis(title="Frequency (number of words)")
		)
		fig = Figure(data=idata, layout=layout)
		if basename != "ALL":
			type_map[basename] = bar

		print("generating plot: {}".format(basename))
		
		plot_url_map[basename] = py.plot(fig, filename='{}-distribution'.format(basename), auto_open=False)

layout = Layout(
	title="Combined Length Distributions (all types)",
	xaxis=XAxis(title="Word length (characters)"),
	yaxis=YAxis(title="Frequency (number of words)"),
	barmode='stacked'
)
fig = Figure(data=type_map.values(), layout=layout)

print("generating stacked plot")

stacked_url = py.plot(fig, filename='combined_word_distributions', auto_open=False)

with open("plot_urls.txt", "w") as f:
	for k in plot_url_map:
		if k != "ALL":
			f.write("{}\n".format(plot_url_map[k]))
	f.write("{}\n".format(plot_url_map["ALL"]))
	f.write(stacked_url)
