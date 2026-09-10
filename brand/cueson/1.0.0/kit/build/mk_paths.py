#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Expand the Gate 1-approved Cueson centerlines with glyphkit primitives."""

from __future__ import annotations

import json
import os
import sys


TEMPLATES = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "..", "..", "skill", "templates"))
if TEMPLATES not in sys.path:
    sys.path.insert(0, TEMPLATES)

import glyphkit as G


GRID = 512

# These centerlines and widths are the approved cueframe-r1 construction.
# Each line is expanded without reinterpretation through glyphkit.capsule.
FULL_BRACKETS = (
    ((168, 104), (108, 104), (108, 408), (168, 408)),
    ((344, 104), (404, 104), (404, 408), (344, 408)),
)
FULL_BRACKET_WIDTH = 34
FULL_RAILS = (
    ((184, 172), (336, 172)),
    ((184, 256), (304, 256)),
    ((184, 340), (352, 340)),
)
FULL_RAIL_WIDTH = 32

REDUCED_BRACKETS = (
    ((174, 126), (118, 126), (118, 386), (174, 386)),
    ((338, 126), (394, 126), (394, 386), (338, 386)),
)
REDUCED_BRACKET_WIDTH = 42
REDUCED_RAILS = (((190, 256), (322, 256)),)
REDUCED_RAIL_WIDTH = 42


def expand_segments(role, chains, width):
    """Expand every consecutive centerline segment as a round-ended fill."""
    paths = []
    for chain in chains:
        for start, end in zip(chain, chain[1:]):
            paths.append({"role": role, "d": G.capsule(start[0], start[1], end[0], end[1], width)})
    return paths


full = (
    expand_segments("accent", FULL_BRACKETS, FULL_BRACKET_WIDTH)
    + expand_segments("ink", FULL_RAILS, FULL_RAIL_WIDTH)
)

reduced = (
    expand_segments("accent", REDUCED_BRACKETS, REDUCED_BRACKET_WIDTH)
    + expand_segments("ink", REDUCED_RAILS, REDUCED_RAIL_WIDTH)
)


if __name__ == "__main__":
    print(json.dumps({"full": full, "reduced": reduced}, indent=2))
